import React from 'react';
import { Modal, Card, Icon } from 'antd';
import styles from './index.module.css';

class AddSelect extends React.Component {
  constructor(props) {
    super(props);
  }

  switchExt1 = () => {
    this.props.ext1Visible
  };

  switchExt2 = () => {
    this.props.ext2Visible
  };

  render() {
    const modalStyle = {
      display: 'flex',
      justifyContent: 'space-around',
      backgroundColor: 'rgba(240, 242, 245, 1)',
      padding: '80px 0'
    };

    return (
      <Modal
        visible
        width={800}
        maskClosable={false}
        title="选择发布方式"
        bodyStyle={modalStyle}
        onCancel={this.props.cancelAddVisible}
        footer={null}>
        <Card
          style={{width: 300, cursor: 'pointer'}}
          bodyStyle={{display: 'flex'}}
          onClick={this.props.ext1Visible}>
          <div style={{marginRight: 16}}>
            <Icon type="ordered-list" style={{fontSize: 36, color: '#1890ff'}}/>
          </div>
          <div>
            <div className={styles.cardTitle}>常规发布</div>
            <div className={styles.cardDesc}>
              由发布平台 来控制发布的主流程，你可以通过添加钩子脚本来执行额外的自定义操作。
            </div>
          </div>
        </Card>
        <Card
          style={{width: 300, cursor: 'pointer'}}
          bodyStyle={{display: 'flex'}}
          onClick={this.props.ext2Visible}>
          <div style={{marginRight: 16}}>
            <Icon type="build" style={{fontSize: 36, color: '#1890ff'}}/>
          </div>
          <div>
            <div className={styles.cardTitle}>自定义发布</div>
            <div className={styles.cardDesc}>
              你可以完全自己定义发布的所有流程和操作，发布平台 负责按顺序依次执行你记录的动作。
            </div>
          </div>
        </Card>
      </Modal>
    )
  }
}

export default AddSelect
