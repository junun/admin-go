import React from 'react';
import styles from './index.module.css';
import lds from 'lodash';

class OutView extends React.Component {
  constructor(props) {
    super(props);
    this.el = null;
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.el != null) {
      this.el.scrollTop = this.el.scrollHeight
    }
  }

  render() {
    const outputs = lds.get(this.props.outputs, `${this.props.id}.Data`, []);
    return (
      <pre ref={el => this.el = el} className={styles.ext1Console}>
        {outputs}
      </pre>
    )
  }
}

export default OutView