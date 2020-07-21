import React, {Fragment, Component} from 'react';
import { Link } from 'react-router-dom';
import { Modal, Form, Select, Button, Icon, Input, Col, Steps } from 'antd';
import styles from './index.module.css';
import Template from './Template';
import Ex1Info from './Ex1Info';
import Ex2Info from './Ex2Info';
import lds from 'lodash';

const FormItem = Form.Item;
const Option = Select.Option;

class SelectTemplate extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      page: 0,
      extend: 1,
      info:{},
    }
  }

  componentDidMount() {
    this.setState({
      info: this.props.editCacheData,
    })
  }

  checkExtend = (tid) => {
    for (var i = 0; i < this.props.appTemplateList.length; i++) {
      if (this.props.appTemplateList[i].Dtid === tid) {
        this.setState({
          extend: this.props.appTemplateList[i].Extend,
        })
        break;
      }
    }
  }

  handler = (values) => {
    this.setState({
      info: values,
      extend: values.Extend,
      page: this.state.page + 1,
    })
    
  }

  prehandler = (values) => {
    this.setState({
      page: this.state.page - 1,
    })
  }

  handlerCancle = () => {
    // return;
    this.props.onCancel();

  }

  render() {
    const {appTemplateList} = this.props;

    return (
      <Modal
        title={Object.keys(this.props.editCacheData).length > 0 ? '项目上线编辑' : '项目上线提单'}
        visible
        width={960}
        maskClosable={false}
        destroyOnClose= "true"
        onCancel={() => this.handlerCancle()}
        footer={null}
      >
        <Steps current={this.state.page} className={styles.steps}>
          <Steps.Step key={0} title="选择发布模板"/>
          <Steps.Step key={1} title="选择发布版本/Tag"/>
        </Steps>
        {this.state.page === 0 && 
          <Template
            appTemplateList={appTemplateList}
            nextPage={this.handler}
            info={this.state.info}
          />
        }
        {this.state.page === 1 &&  this.state.extend === 1 &&
          <Ex1Info
            prePage={this.prehandler}
            info={this.state.info}
            onCancel={this.props.onCancel}
          />
        }
        {this.state.page === 1 &&  this.state.extend === 2 &&
          <Ex2Info
            prePage={this.prehandler}
            info={this.state.info}
            onCancel={this.props.onCancel}
          />
        }
      </Modal>
    )
  }
}

export default SelectTemplate
